import os
import subprocess


class Builder(object):
    def detect_go_binary(self, binary="go.exe", delim=";"):
        paths = os.environ["PATH"].split(delim)
        for path in paths:
            if not os.path.exists(path):
                continue

            files = os.listdir(path)
            files = [file for file in files if file == binary]
            if files:
                return os.path.join(path, binary)

    def run(self):
        src_dir = "src"
        target = os.path.join("..", "bin", "havenet.exe")

        files = os.listdir(src_dir)
        files.sort()
        files = [file for file in files if ".go" in file]
        files = [file for file in files if not "_test.go" in file]

        executable = self.detect_go_binary()
        cwd = "src"
        args = [executable] + ["build"] + ["-o", target] + files

        print("Invoking [%s] %s" % (cwd, args))
        proc = subprocess.Popen(args=args, cwd=cwd)
        proc.communicate()


if __name__ == '__main__':
    Builder().run()
