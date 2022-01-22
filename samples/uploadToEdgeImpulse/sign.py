import argparse
import hashlib
import hmac
import json
import sys


def sign(data, hmac_key):
    file_hmac = hmac.new(bytes(hmac_key, 'utf-8'), msg=data, digestmod=hashlib.sha256).hexdigest()
    file_len = len(data)

    # empty signature (all zeros). HS256 gives 32 byte signature, and we encode in hex, so we need 64 characters here
    empty_signature = ''.join(['0'] * 64)
    data = {
        "protected": {
            "ver": "v1",
            "alg": "HS256",
        },
        "signature": empty_signature,
        "payload": {
            "device_name": "test-dev",
            "device_type": "phone",
            "interval_ms": 0,
            "sensors": [
                {"name": "image", "units": "rgba"}
            ],
            "values": ["Ref-BINARY-image/jpeg (%d bytes) %s" % (file_len, file_hmac)]
        }
    }

    # encode in JSON
    encoded = json.dumps(data)
    # sign message
    signature = hmac.new(bytes(hmac_key, 'utf-8'), msg=encoded.encode('utf-8'), digestmod=hashlib.sha256).hexdigest()
    # set the signature again in the message, and encode again
    data['signature'] = signature
    return json.dumps(data)


def main(argv):
    parser = argparse.ArgumentParser(description="Generates signed metadata file for EdgeImpluse upload",
                                     epilog="sys.argv[0] input.jpg",
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument("--hmac-key", default='', type=str, help="Edge Impulse HMAC Key", required=True)
    parser.add_argument("file", default='', type=str, help="File to sign")
    parser.add_argument("output", default='', type=str, help="Output file")

    args = parser.parse_args(args=argv)

    print("Starting", file=sys.stderr)
    print("In : %s" % args.file, file=sys.stderr)
    print("Out: %s" % args.output, file=sys.stderr)
    
    with open(args.file, "rb") as file:
        data = file.read()
        res = sign(data, args.hmac_key)
        with open(args.output, 'w') as out:
            out.write(res)

    print("Done", file=sys.stderr)


if __name__ == "__main__":
    main(sys.argv[1:])