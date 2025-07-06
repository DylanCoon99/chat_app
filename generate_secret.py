import os 

with open(".env", "a") as f:
	f.write(os.urandom(32).hex())