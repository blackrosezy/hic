# -*- coding: utf-8 -*-
from distutils.core import setup

from setuptools import find_packages


setup(
    name='hic',
    version='2.0.0',
    author='Mohd Rozi',
    author_email='blackrosezy@gmail.com',
    url='https://bitbucket.org/blackrosezy/hic',
    license='MIT, see LICENSE file',
    description='Hipache cli tool.',
    long_description=open('README.rst').read(),
    install_requires=[
        'docopt==0.6.2',
        'docker-py==0.5.3',
        'redis==2.10.3',
        'tabulate==0.7.3'
    ],
    scripts=['hic'],
    zip_safe=False,
)