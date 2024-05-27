#!/usr/bin/env python
# -*- coding: utf-8 -*-
import time
from locust import task,TaskSet,FastHttpUser, between

USER_CREDENTIALS = [
    {
        "username": f"user_{i}",
        "password": f"psw_{i}"
    }
    for i in range(1, 301)
]

class WebsiteTasks(TaskSet):
    @task
    def login(self):
        # Define the login payload
        payload = {
            "username":  self.username,
            "password":  self.password
        }

        # Send a POST request to the login endpoint
        headers = {'Content-Type': 'application/json'}
        self.client.post("/api/users/login", json=payload, headers=headers)

    def on_start(self):
        if len(USER_CREDENTIALS) > 0:
            # Pop a set of credentials for the user from the list
            self.username, self.password = USER_CREDENTIALS.pop()
        else:
            # If no more credentials are available, stop this user
            self.interrupt(reschedule=False)

class WebsiteUser(FastHttpUser):
    tasks = [WebsiteTasks]
    host = "http://localhost:8080"
    min_wait = 1000
    max_wait = 2000

import subprocess

def start_locust_master():
    return subprocess.Popen(["locust", "-f", "/Users/zihehuang/GolandProjects/entrytask/test/locustfile.py", "--master", "--host=http://localhost:8080"])

def start_locust_worker():
    return subprocess.Popen(["locust", "-f", "/Users/zihehuang/GolandProjects/entrytask/test/locustfile.py", "--worker", "--master-host=127.0.0.1"])


if __name__ == "__main__":
    subprocessesList = []
    subprocessesList.append(start_locust_master())
    time.sleep(3)  # Wait for the master to start
    num = 8
    for i in range(1, num):
        subprocessesList.append(start_locust_worker())
    # Start subprocesses

    # Wait for subprocesses to complete
    input("Press Enter to continue...")
    print("Continuing the script...")

    # Terminate subprocesses
    for process in subprocessesList:
        process.terminate()