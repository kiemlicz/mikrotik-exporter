from setuptools import setup, find_packages, Command
from os import path
import subprocess


class BuildDockerImageCommand(Command):
    description = "Build a Docker image using the provided Dockerfile"
    user_options = []

    def initialize_options(self):
        pass

    def finalize_options(self):
        pass

    def run(self):
        subprocess.check_call(["docker", "build", "-t", "mikrotik_exporter:latest", "."])


setup(
    name="mikrotik_exporter",
    version="0.0.1",
    install_requires=[
        'prometheus-client>=0.21.1',
        'RouterOS-api>=0.21.0',
        'dynaconf>=3.2.10'
    ],
    entry_points={
        'console_scripts': [
            'mikrotik-exporter=exporter.cli.mex:main',
        ],
    },
    cmdclass={
        'build_docker': BuildDockerImageCommand,
    }
)
