from google.cloud import pubsub_v1
from typing import Tuple


def get_instagram_caption(client: pubsub_v1.subscriber.Client, subscription_path: str) -> Tuple[bytes, str]:
    result: bytes = bytes(0)
    try:
        pull_result = client.pull(subscription=subscription_path, max_messages=1, )
    except Exception:
        return result, ""

    if len(pull_result.received_messages) <= 1:
        return result, ""

    result = pull_result.received_messages[0].message.data
    return result, pull_result.received_messages[0].ack_id
