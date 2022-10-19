from celery import Celery
import requests
app = Celery(
    'tasks',
    broker='redis://',
    backend='redis://',
)

app.conf.update(
    CELERY_TASK_SERIALIZER='json',
    CELERY_ACCEPT_CONTENT=['json'],  # Ignore other content
    CELERY_RESULT_SERIALIZER='json',
    CELERY_ENABLE_UTC=True,
    CELERY_TASK_PROTOCOL=1,
)

@app.task
def postRequest(message, address):
    return requests.post(address, message)

@app.task
def sendMail(message, emailAddress, seconds):
    postRequest.apply_async(({"message": message}, emailAddress),countdown=seconds)
