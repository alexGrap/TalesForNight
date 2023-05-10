#!/usr/bin/python3
import os
import codecs
import pyttsx3

engine = pyttsx3.init()

engine.setProperty('rate', 150)
engine.setProperty('volume', 0.7)


fileObj = codecs.open( "../temp-folder/text.txt", "r", "utf_8_sig" )
text = fileObj.read()

engine.save_to_file(text, '../temp-folder/file.ogg')
fileObj.close()
os.remove('../temp-folder/text.txt')
engine.runAndWait()