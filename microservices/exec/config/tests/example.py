import os
import sys
import shutil
import time
import random

print(sys.argv)
infile = sys.argv[1]
print(infile)
output_dir = sys.argv[2]
print(output_dir)

print("process start")
time.sleep(random.randint(0, 5))
print("process end")

# raise BaseException("python is error")

print("move start")
outfile = os.path.join(output_dir, os.path.basename(infile))
shutil.copy(infile, outfile + ".json1")
shutil.copy(infile, outfile + ".json2")
print("move end")
