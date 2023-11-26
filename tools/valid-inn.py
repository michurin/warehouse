#!/usr/bin/env python3

'''
Для 12-ти значного ИНН генерирует контрольную сумму.
Использование: указываете 10 любых цифр и получаете
полный корректный ИНН:

$ ./valid-inn.py 7777111111
777711111191
'''

import sys

w = [3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8]

inn = sys.argv[1] # must be 10 digits
inn += str(((sum(int(x[0])*x[1] for x in zip(inn, w[1:])))%11)%10)
inn += str(((sum(int(x[0])*x[1] for x in zip(inn, w)))%11)%10)
print(inn)
