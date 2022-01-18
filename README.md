
apt-get install python-opencv python-numpy


git clone https://github.com/groupgets/pylepton.git
git checkout lepton3-dev
python setup.py install

// Without this edit it returns an error OSError: [Errno 90] Message too long
nano pylepton/Lepton.py
SPIDEV_MESSAGE_LIMIT = 8 