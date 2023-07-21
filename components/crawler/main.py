import base64
import os

import instagrapi
import instagrapi.exceptions
import requests
from google.cloud import pubsub_v1
from google.cloud import storage
import gcs
import pubsub
import logging


def main():
    logger = logging.getLogger()
    try:
        project_id = requests.get("http://metadata.google.internal/computeMetadata/v1/project/project-id", headers={
            "Metadata-Flavor": "Google"
        }).text
    except Exception as e:
        project_id = os.getenv("PROJECT_ID")
        if project_id == "":
            logger.error(e)
            return

    notify_endpoint = os.getenv("NOTIFY_ENDPOINT")
    url = os.getenv("INSTAGRAM_POST_CODE")
    ig_username = os.getenv("IG_USERNAME")
    ig_password = os.getenv("IG_PASSWORD")
    subscription_path = os.getenv("SUBSCRIPTION_PATH")
    topic = os.getenv("TOPIC")
    api_key = os.getenv("API_KEY")

    storage_client = storage.Client(project=project_id)
    subscriber = pubsub_v1.SubscriberClient()
    publisher = pubsub_v1.PublisherClient()
    instagrapi_config = gcs.get_instagrapi_config(storage_client)
    ig_client = instagrapi.Client(settings=instagrapi_config)

    old_caption, ack_id = pubsub.get_instagram_caption(subscriber, subscription_path)

    try:
        ig_client.login(ig_username, ig_password)
        ig_pk = ig_client.media_pk_from_code(url)
        detail = ig_client.media_info(ig_pk)
    except instagrapi.exceptions.ClientError as e:
        logger.error(e)
        return

    new_caption = base64.b64encode(detail.caption_text.encode())

    if new_caption != old_caption:
        notify_response = requests.post(f"{notify_endpoint}/notification/line",
                                        json={
                                            "icon_url": detail.user.profile_pic_url,
                                            "message": detail.caption_text,
                                            "name": detail.user.full_name
                                        },
                                        headers={
                                            "X-API-KEY": api_key
                                        }
                                        )

        discord_notify = requests.post(f"{notify_endpoint}/notification/discord",
                                       json=
                                       {
                                           "icon_url": detail.user.profile_pic_url,
                                           "message": detail.caption_text,
                                           "name": detail.user.full_name,
                                           "profile_url": f"https://www.instagram.com/{detail.user.username}",
                                           "title": "PRSMINライブスケジュール",
                                           "post_url": f"https://www.instagram.com/p/{detail.code}"
                                       },
                                       headers={
                                           "X-API-KEY": api_key
                                       }
                                       )
        try:
            notify_response.raise_for_status()
            discord_notify.raise_for_status()
        except Exception as e:
            logger.error(e)
            return

    publish_future = publisher.publish(topic=topic, data=new_caption)

    try:
        publish_future.result()
    except Exception as e:
        logger.error(e)
    else:
        if ack_id != "":
            subscriber.acknowledge(subscription=subscription_path, ack_ids=[ack_id])


if __name__ == "__main__":
    main()
