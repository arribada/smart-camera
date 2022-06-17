
apt-get install python-opencv python-numpy

pip install opencv-python-headless

apt-get install libatlas-base-dev
apt-get install libopenjp2-7
apt-get install libavcodec-dev libavformat-dev libswscale-dev libv4l-dev


git clone https://github.com/groupgets/pylepton.git
git checkout lepton3-dev
python setup.py install

// Without this edit it returns an error OSError: [Errno 90] Message too long
nano pylepton/Lepton.py
SPIDEV_MESSAGE_LIMIT = 8 