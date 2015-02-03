from distutils.core import setup

setup(
    name='hic',
    version='2.0.0',
    packages=['hic'],
    url='https://github.com/blackrosezy/hic',
    license='MIT, see LICENSE file',
    author='MohdRozi',
    author_email='blackrosezy@gmail.com',
    description='Hipache cli tool',
    long_description=open('README.md').read(),
    install_requires=[
        'docopt==0.6.2',
        'docker-py==0.5.3',
        'redis==2.10.3',
        'tabulate==0.7.3'
    ],
    entry_points={
        'console_scripts': ['hic=hic.hic:main'],
    },
    zip_safe=False
)