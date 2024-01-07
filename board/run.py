import pyaudio
import wave
from maix import event
from select import select
import time
import requests

def check_key():
	import os
	tmp = "/dev/input/by-path/"
	if os.path.exists(tmp):
		for i in os.listdir(tmp):
			if i.find("kbd") != -1:
				return tmp + i
	return "/dev/input/event0"

key_released = False  # Boolean variable to track key release
dev = event.InputDevice(check_key())


def record():
    import pyaudio
    import wave
    from maix import event
    from select import select
    import time
    import requests
    
    CHUNK = 1024
    FORMAT = pyaudio.paInt16
    CHANNELS = 2
    RATE = 44100
    RECORD_SECONDS = 2
    WAVE_OUTPUT_FILENAME = "test.wav"

    url = "http://192.168.124.21:8080/data"
    
    p = pyaudio.PyAudio()

    stream = p.open(format=FORMAT,
                    channels=CHANNELS,
                    rate=RATE,
                    input=True,
                    frames_per_buffer=CHUNK)

    print("* recording")

    frames = []

    for i in range(0, int(RATE / CHUNK * RECORD_SECONDS)):
        data = stream.read(CHUNK)
        frames.append(data)

    print("* done recording")

    stream.stop_stream()
    stream.close()
    p.terminate()

    wf = wave.open(WAVE_OUTPUT_FILENAME, 'wb')
    wf.setnchannels(CHANNELS)
    wf.setsampwidth(p.get_sample_size(FORMAT))
    wf.setframerate(RATE)
    wf.writeframes(b''.join(frames))
    wf.close()
    
    files = {"file": open(WAVE_OUTPUT_FILENAME, "rb")}
    response = requests.post(url, files=files)
    
    if response.status_code == requests.codes.ok:
        file_name = "response.wav"

        with open(file_name, "wb") as file:
            file.write(response.content)
    
        print("File downloaded and saved successfully.")
    else:
        print("File download failed.")
    
    # play
    import wave
    import sys
    wf = wave.open(r'response.wav', 'rb')#(sys.argv[1], 'rb'
    p = pyaudio.PyAudio()

    stream = p.open(format=p.get_format_from_width(wf.getsampwidth()),
                    channels=wf.getnchannels(),
                    rate=wf.getframerate(),
                    output=True)

    data = wf.readframes(CHUNK)

    while len(data) > 0:
        stream.write(data)
        data = wf.readframes(CHUNK)

    stream.stop_stream()
    stream.close()

    p.terminate()

while True:
    r, w, x = select([dev], [], [], 0) # If r == 0 or set to 0, read() will raise a BlockingIOError
    if r:
        for data in dev.read():
            print(data)
            if data.code == 0x03:
                if data.value == 1:  # Key press event
                    key_released = False
                    print('Pressed key S2')
                    record()
                    
    else:  # Key release event
        if not key_released:
            key_released = True
            print('Released key S2')
    
    time.sleep(1) # Sleep for 1 second
