import arvados
import arvados_testutil as tutil
import hashlib

class ManifestExamples(object):
    def make_manifest(self,
                      bytes_per_block=1,
                      blocks_per_file=1,
                      files_per_stream=1,
                      streams=1):
        datablip = 'x' * bytes_per_block
        data_loc = '{}+{}'.format(hashlib.md5(datablip).hexdigest(),
                                  bytes_per_block)
        with tutil.mock_keep_responses(data_loc, 200):
            coll = arvados.CollectionWriter()
            for si in range(0, streams):
                for fi in range(0, files_per_stream):
                    with coll.open("stream{}/file{}.txt".format(si, fi)) as f:
                        for bi in range(0, blocks_per_file):
                            f.write(datablip)
            return coll.manifest_text()