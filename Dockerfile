FROM python:2.7-slim-buster

#docker run -it --privileged -v /home:/home/host python:2.7-slim-buster /bin/bash

ENV PYTHONPATH=/usr/lib/python2.7/dist-packages
RUN apt-get update && apt-get install -y --no-install-recommends git python-opencv python-numpy
RUN git clone https://github.com/groupgets/pylepton.git
RUN cd pylepton/
RUN git checkout lepton3-dev
RUN python setup.py install