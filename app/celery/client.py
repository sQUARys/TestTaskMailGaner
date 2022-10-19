from worker import add, add_reflect
from datetime import timedelta,datetime

ar = add_reflect.apply_async(kwargs={'a': 5456, 'b': 2878}, serializer='json', expires=120)
print('Result: {}'.format(ar.get()))