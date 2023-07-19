FROM python:3.11-slim

WORKDIR /usr/src/app

COPY components/crawler/requirements.txt /usr/src/app
RUN pip install --no-cache-dir -r requirements.txt
COPY components/crawler/ /usr/src/app/

CMD [ "python", "main.py" ]