import json
import os

import google.cloud.exceptions
from google.cloud import storage


def get_instagrapi_config(client: storage.client.Client) -> dict:
    bucket_name = os.getenv("GCS_INSTAGRAPI_CONFIG_BUCKET", default="")
    blob_key = os.getenv("GCS_INSTAGRAPI_CONFIG_KEY", default="")

    if bucket_name or blob_key == "":
        return {}

    try:
        config_json: str = client.get_bucket(bucket_name).get_blob(blob_key).download_as_text()
    except google.cloud.exceptions.NotFound:
        return {}

    config = json.loads(config_json)

    return config
