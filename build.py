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
    binary_name = '%s.exe' % binary if os_is_windows() else binary

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
    sys.stdout.write(msg)
    sys.stdout.flush()


class Builder(object):
    version_module = os.path.join('src', 'app_version.go')

    def detect_version(self):
        git_exe = detect_binary(binary='git')
        if git_exe is None:
            return '?'

        args = [git_exe, 'describe']
        version = invoke(args=args)
        return version.decode('ascii')

    def set_version(self, version):
        version = version or '?'
        with open(self.version_module, 'wb') as f:
            content = 'package main\n\n\nconst appVersion = "%s"\n' % version
            content_b = content.encode('ascii')
            f.write(content_b)

    def build(self):
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

    def run(self, force_build=False):
        version = self.detect_version()
        self.set_version(version)

        # If we're forcing or on Windows build here, otherwise let the makefile
        # handle it
        if force_build or os_is_windows():
            self.build()


if __name__ == '__main__':
    from optparse import OptionParser

    parser = OptionParser()
    parser.add_option('', '--force', action='store_true',
                      help='Always build using build.py')
    (options, args) = parser.parse_args()

    builder = Builder()
    builder.run(force_build=options.force)
