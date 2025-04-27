import routeros_api


class ROS:
    def __init__(self, host: str, username: str, password: str):
        self.client = routeros_api.RouterOsApiPool(
            host,
            username=username,
            password=password,
            port=8728,
            plaintext_login=True,
        )
        self.connection = None

    def __enter__(self):
        self.connection = self.client.get_api()
        return self

    def __exit__(self, exception_type, exception_value, exception_traceback):
        try:
            self.client.disconnect()
        except Exception as e:
            ???
