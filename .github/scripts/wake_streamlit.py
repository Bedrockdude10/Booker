"""
Wake up a sleeping Streamlit Community Cloud app.
Clicks the "Yes, get this app back up!" button if present.
"""

import os
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.chrome.service import Service


def wake_streamlit_app(url: str) -> None:
    options = Options()
    options.add_argument("--headless")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")

    driver = webdriver.Chrome(
        service=Service(ChromeDriverManager().install()),
        options=options
    )

    try:
        print(f"Visiting {url}...")
        driver.get(url)

        # Check if the app is sleeping (wake button present)
        try:
            wake_button = WebDriverWait(driver, 15).until(
                EC.element_to_be_clickable((By.XPATH, "//*[contains(text(), 'Yes, get this app back up')]"))
            )
            print("App is sleeping — clicking wake button...")
            wake_button.click()
            time.sleep(10)  # Wait for app to start booting
            print("Wake button clicked. App is booting.")
        except TimeoutException:
            print("App is already awake — no action needed.")

    finally:
        driver.quit()


if __name__ == "__main__":
    url = os.environ.get("STREAMLIT_URL")
    if not url:
        raise ValueError("STREAMLIT_URL environment variable not set")
    wake_streamlit_app(url)