'''This needs to be a script as passing *.go to "go build" in the terminal on
Windows does not work.'''


import os
import subprocess
import sys


def detect_binary_once(binary_name):
    '''Windows requires a full path, so detect the full path of a binary by
    traversing $PATH.'''

    delim = ';' if os_is_windows() else ':'

    paths = os.environ["PATH"].split(delim)
    for path in paths:
        if not os.path.exists(path):
            continue

        files = os.listdir(path)
        files = [file for file in files if file == binary_name]
        if files:
            return os.path.join(path, binary_name)

def detect_binary(binary):
    binary_name = '%s.exe' if os_is_windows() else binary

    filepath = detect_binary_once(binary_name)

    # Try again without .exe extension
    if filepath is None and os_is_windows():
        filepath = detect_binary_once(binary)

    return filepath


def invoke(args, cwd='.'):
    write("Invoking [cwd: %s] %s" % (cwd, args))
    proc = subprocess.Popen(
        args=args, cwd=cwd,
        stdout=subprocess.PIPE, stderr=subprocess.PIPE,
    )
    stdout, stderr = proc.communicate()
    return stdout.strip()

def os_is_windows():
    return sys.platform.startswith('win')

def write(msg):
    msg = '%s\n' % msg
    line = msg.encode('ascii')
    sys.stdout.write(line)
    sys.stdout.flush()


class Builder(object):
    version_module = os.path.join('src', 'app_version.go')

    def detect_version(self):
        git_exe = detect_binary(binary='git')
        if git_exe is None:
            return '?'

        args = [git_exe, 'describe']
        version = invoke(args=args)
        return version

    def set_version(self, version):
        version = version or '?'
        with open(self.version_module, 'wb') as f:
            content = 'package main\n\n\nconst appVersion = "%s"\n' % version
            content_b = content.encode('ascii')
            f.write(content_b)

    def build_on_windows(self):
        src_dir = "src"
        target = os.path.join("..", "bin", "havenet.exe")

        # Find all the sources
        files = os.listdir(src_dir)
        files.sort()
        files = [file for file in files if ".go" in file]
        files = [file for file in files if not "_test.go" in file]

        cwd = "src"
        executable = detect_binary(binary='go')
        args = [executable] + ["build"] + ["-o", target] + files

        invoke(args=args, cwd=cwd)

    def run(self):
        version = self.detect_version()
        self.set_version(version)

        # If we're not on windows we let the makefile handle the build
        if not os_is_windows():
            return

        self.build_on_windows()


if __name__ == '__main__':
    builder = Builder()
    builder.run()
