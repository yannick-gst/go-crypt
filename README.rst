=======
gocrypt
=======

**gocrypt** is a command line tool that encrypts and decrypts files using the Advanced Encryption Standard (AES). It is written in Go.

Building
========
To build, simply run ``go build``.

Testing
=======
To run tests, simply run ``go test``.

Usage
=====
Running ``gocrypt --help`` or ``goccrypt -h`` will display help information:

.. code-block:: console

    usage: gocrypt [-h|--help] -i|--input "<value>" -o|--output "<value>" -m|--mode
               (encrypt|decrypt)

               A command line tool to encrypt a file, or decrypt a file that
               was encrypted with this tool.

    Arguments:

    -h  --help    Print help information
    -i  --input   Path to input file
    -o  --output  Path to output file
    -m  --mode    Specifies whether to encrypt or decrypt

To encrypt a file, one can run:

.. code-block:: console

    gocrypt -i file_to_encrypt -o output_file -m encrypt

To decrypt a file, one can run:

.. code-block:: console

    gocrypt -i file_to_decrypt -o output_file -m decrypt

For both encryption and decryption, the user will be prompted to enter a password. The password used for decryption should match the password that was used to encrypt the original data.

License
=======
This project is licensed under the MIT license. See the `LICENSE <https://github.com/yannick-gst/gocrypt/blob/main/LICENSE>`_ file for details.

Third-Party Libraries
=====================
This project uses `package pbkdf2 <https://pkg.go.dev/golang.org/x/crypto/pbkdf2>`_ for password-based file encryption.
