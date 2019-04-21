#!/usr/bin/env liquidsoap


%include "utils.liq"


# CONFIGS
set("log.file", true)
set("log.stdout", true)
set("log.level", 4)
set("server.telnet", true)
set("server.telnet.port", 5000)
set("log.file.path", "/tmp/broadcast.log")
set("audio.converter.samplerate.libsamplerate.quality","best")
set("buffering.path","/tmp")
set("decoding.buffer_length",10.)


# FUNCTIONS
output_docker_env.icecast = output.icecast(
    description="Radio Imaginee de la Suite du Monde",
    genre="Documentaires",
    url="https://www.lasuitedumonde.com/",
    host=my_getenv("ICECAST_PORT_8000_TCP_ADDR"),
    port=int_of_string(my_getenv("ICECAST_PORT_8000_TCP_PORT")),
    password=my_getenv("ICECAST_SOURCE_PASSWORD")
)

def outputs(encoder, ~mount, input)
    output_docker_env.icecast(encoder, mount=mount, input)
end

# MIXS
default = mksafe(
    rewrite_metadata(
        [
            ("artist", "La Suite du Monde"),
            ("title", "Maintenance en cours"),
            ("comment", "https://www.lasuitedumonde.com/")
            ],
        audio_to_stereo(
            playlist.safe(
                reload=60,
                "/playlists/failures/",
                mode="random"
                )
            )
        )
    )

harbor_input = input.harbor("main.harbor", port=5001, password=my_getenv("HARBOR_PASSWORD"))

input = mksafe(fallback(track_sensitive=false, [ harbor_input, default ]))


# NORMALIZE
#input = audio_process(input)


# OUTPUTS
outputs(%mp3(bitrate=64), mount="/radio-imaginee-64.mp3", input)
outputs(%mp3(bitrate=128), mount="/radio-imaginee-128.mp3", input)
outputs(%mp3(bitrate=160), mount="/radio-imaginee-160.mp3", input)
outputs(%mp3(bitrate=192), mount="/radio-imaginee-192.mp3", input)
outputs(%aac(bitrate=128), mount="/radio-imaginee-128.aac", input)
outputs(%aac(bitrate=160), mount="/radio-imaginee-160.aac", input)
outputs(%aac(bitrate=192), mount="/radio-imaginee-192.aac", input)
outputs(%vorbis, mount="/radio-imaginee.ogg", input)
