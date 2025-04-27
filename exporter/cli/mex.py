import argparse

from dynaconf import Dynaconf
from prometheus_client import start_http_server

parser = argparse.ArgumentParser()
parser.add_argument('--app_host', help='Override app host')
args = parser.parse_args()


def config() -> Dynaconf:
    settings = Dynaconf(
        settings_files=['settings.yaml'],
        envvar_prefix='MEX',
    )
    return settings


def main() -> None:
    c = config()
    server, t = start_http_server(c.metrics.port, addr=c.metrics.host)
    # todo handle lifecycle
    # server.shutdown()
    # t.join()


if __name__ == '__main__':
    main()
