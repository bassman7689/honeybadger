import smtplib
import sys

port = 2525

try:
    server = smtplib.SMTP("localhost", port)
    server.ehlo()
    server.sendmail("test@test.com", "test@test.com", """Subject: Hi there

            This message is sent from Python.""")
except:
    e = sys.exc_info()[0]
    print(e)

