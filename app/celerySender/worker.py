from celery import Celery
import requests
app = Celery(
    'tasks',
    broker='redis://',
    backend='redis://',
)

app.conf.update(
    CELERY_TASK_SERIALIZER='json',
    CELERY_ACCEPT_CONTENT=['json'],
    CELERY_RESULT_SERIALIZER='json',
    CELERY_ENABLE_UTC=True,
    CELERY_TASK_PROTOCOL=1,
)

@app.task
def postRequest(html, address):
    requests.post(address, data=html)


@app.task
def sendMail(html, emailAddress, seconds):
    postRequest.apply_async((html, emailAddress), countdown=seconds)