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
        logging.error(e)
        return

    notify_endpoint = os.getenv("NOTIFY_ENDPOINT")
    url = os.getenv("INSTAGRAM_POST_CODE")
    ig_username = os.getenv("IG_USERNAME")
    ig_password = os.getenv("IG_PASSWORD")
    subscription_path = os.getenv("SUBSCRIPTION_PATH")
    topic = os.getenv("TOPIC")

    storage_client = storage.Client(project=project_id)
    subscriber = pubsub_v1.SubscriberClient(project=project_id)
    publisher = pubsub_v1.PublisherClient(project=project_id)
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

    new_caption = base64.b64encode(detail.caption_text)

    if new_caption != old_caption:
        requests.post(notify_endpoint, json={
            "icon_url": detail.user.profile_pic_url_hd,
            "message": detail.caption_text,
            "name": detail.user.full_name
        })

    try:
        publisher.publish(topic=topic, data=new_caption)
    except:
        return
    else:
        subscriber.acknowledge(ack_ids=[ack_id])


if __name__ == "__main__":
    main()
