from distutils.core import setup

setup(
    name='hic',
    version='2.0.0',
    packages=['hic'],
    url='https://github.com/blackrosezy/hic',
    license='MIT',
    author='MohdRozi',
    author_email='blackrosezy@gmail.com',
    description='Hipache cli tool',
    entry_points={
        'console_scripts': ['hic=hic.hic:main'],
    }
)
