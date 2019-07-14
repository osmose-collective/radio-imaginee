#!/usr/bin/env liquidsoap


%include "utils.liq"


# CONFIGS
set("log.file", true)
set("log.stdout", true)
set("log.level", 3)
set("server.telnet", true)
set("server.telnet.port", 5000)
set("server.telnet.bind_addr", "0.0.0.0")
set("log.file.path", "/tmp/broadcast.log")
set("audio.converter.samplerate.libsamplerate.quality","best")
set("buffering.path","/tmp")
set("decoding.buffer_length",10.)


# INPUTS
#clock = mksafe(
#    audio_to_stereo(
#        playlist.safe(
#            reload=3600,
#            "/playlists/clock/",
#            mode="randomize",
#            reload_mode="rounds"
#        )
#    )
#)

def write_metadata_to_file(m) =
    system("echo "^ quote(m["artist"] ^" - "^ m["title"]) ^" > /data/latest.txt")
    system("echo "^ quote(m["artist"] ^" - "^ m["title"]) ^" >> /data/history.txt")
end

jingles = mksafe(
    rewrite_metadata(
        [
            ("artist", "Radio imaginee"),
            ("title","Jingle"),
            ("comment","https://radio.osmose.world/")
            ],
        audio_to_stereo(
            playlist.safe(
                reload=600,
                "/playlists/jingles",
                mode="randomize",reload_mode="rounds"
                )
            )
        )
    )

songs = mksafe(
    audio_to_stereo(
        rewrite_metadata(
            [
                ("title", "$(title) (Radio imaginee)"),
                ("comment", "https://radio.osmose.world/")
                ],
            playlist.safe(
                reload=600,
                "/playlists/musiques",
                mode="randomize",reload_mode="rounds"
                )
            )
        )
    )

emissions = mksafe(
    audio_to_stereo(
        rewrite_metadata(
            [
#                ("title", "$(title) (emissions - sdm)"),
                ("title", "Emission imaginee"),
                ("comment", "https://radio.osmose.world/")
                ],
            playlist.safe(
                reload=600,
                "/playlists/emissions",
                mode="randomize",reload_mode="rounds"
                )
            )
        )
    )

live_priority = audio_to_stereo(
    rewrite_metadata(
        [
            ("title", "$(title) (LIVE)"),
            ("comment", "https://radio.osmose.world/")
           ],
        stretch(
            ratio=interactive.float("test", 1.00005),
            input.harbor(
                "priority.live",
                id="priority.live",
                port=5002,
                password=my_getenv("LIVE_PASSWORD"),
                user="source",
                logfile="/tmp/harbor.log"
                )
            )
        )
    )

live_guest = audio_to_stereo(
    rewrite_metadata(
        [
            ("title", "$(title) (LIVE)"),
            ("comment", "https://radio.osmose.world/"),
            ],
        stretch(
            ratio=interactive.float("test", 1.00005),
            input.harbor(
                "normal.live",
                id="normal.live",
                port=5001,
                password=my_getenv("LIVE_PASSWORD"),
                user="source",
                logfile="/tmp/harbor.log"
                )
            )
        )
    )


# MIXS
#calendar_mix = switch(track_sensitive=true, [
#        ({ 23h-23h59m59s and 4w }, songs_nuit ),
#        ({ 0h-6h and 5w }, songs_nuit ),
#        ({ 23h-23h59m59s and 5w }, songs_nuit ),
#        ({ 0h-6h and 6w }, songs_nuit ),
#        ({ 23h-23h59m59s and 6w }, songs_nuit ),
#        ({ 0h-6h and 7w }, songs_nuit ),
#        ({ true }, songs_journee)
#])

emissions_mix = switch([
	 ({9h and 0m and 0s-5s},emissions),
	 ({12h and 0m and 0s-5s},emissions),
	 ({15h and 0m and 0s-5s},emissions),
	 ({18h and 0m and 0s-5s},emissions),
	 ({22h and 0m and 0s-5s},emissions),
])

auto_mix = songs # FIXME: use a genre-based scheduler
#auto_mix = dole_fade(jingles=jingles, auto_mix)
auto_mix = fade(jingles, auto_mix)

mix = fallback(track_sensitive=false, [live_priority, live_guest, emissions_mix, auto_mix])

special = request.queue(id="special")
mix = smooth_add(normal=mix,special=special)

# normalizer
mix = audio_process(mix)

# OUTPUTS
mix = mksafe(mix)

mix = on_metadata(write_metadata_to_file, mix)

output.icecast(
    password=my_getenv("HARBOR_PASSWORD"),
    host=my_getenv("BROADCAST_PORT_5001_TCP_ADDR"),
    port=int_of_string(my_getenv("BROADCAST_PORT_5001_TCP_PORT")),
    %wav,
    mount="main.harbor",
    mix
    )

#title = '%Y-%m-%d, %H:%M:%S - $(if $(artist), "$(artist)", "Artiste inconnu") - $(if $(title),"$(title)", "Emission inconnue")'
#output.file(
#    %vorbis,
#     reopen_on_metadata=false,
#     fallible=true,
#     "/playlists/emissions/scc/" ^ title ^ ".ogg",
#     live_scc
#     )
#output.file(
#    %vorbis,
#     reopen_on_metadata=false,
#     fallible=true,
#     "/playlists/emissions/guest/" ^ title ^ ".ogg",
#     live_guest
#     )


