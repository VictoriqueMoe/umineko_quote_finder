import { useCallback, useRef, useState } from "react";

export interface AudioPlayerState {
    isPlaying: boolean;
    currentTime: number;
    duration: number;
    activeId: string | null;
}

export function useAudioPlayer() {
    const audioRef = useRef<HTMLAudioElement | null>(null);
    const [state, setState] = useState<AudioPlayerState>({
        isPlaying: false,
        currentTime: 0,
        duration: 0,
        activeId: null,
    });

    const getAudio = useCallback((): HTMLAudioElement => {
        if (!audioRef.current) {
            const audio = new Audio();
            const savedVolume = localStorage.getItem("uminekoVolume");
            if (savedVolume) {
                audio.volume = parseFloat(savedVolume);
            }
            audio.addEventListener("timeupdate", () => {
                setState(prev => ({
                    ...prev,
                    currentTime: audio.currentTime,
                    duration: audio.duration,
                }));
            });
            audio.addEventListener("loadedmetadata", () => {
                setState(prev => ({
                    ...prev,
                    duration: audio.duration,
                }));
            });
            audio.addEventListener("ended", () => {
                setState(prev => ({
                    ...prev,
                    isPlaying: false,
                    currentTime: 0,
                    activeId: null,
                }));
            });
            audioRef.current = audio;
        }
        return audioRef.current;
    }, []);

    const play = useCallback(
        (url: string, id: string) => {
            const audio = getAudio();

            if (state.activeId === id) {
                if (audio.paused) {
                    audio.play();
                    setState(prev => ({ ...prev, isPlaying: true }));
                } else {
                    audio.pause();
                    setState(prev => ({ ...prev, isPlaying: false }));
                }
                return;
            }

            audio.pause();
            audio.src = url;
            audio.play();
            setState({
                isPlaying: true,
                currentTime: 0,
                duration: 0,
                activeId: id,
            });
        },
        [getAudio, state.activeId],
    );

    const stop = useCallback(() => {
        const audio = audioRef.current;
        if (audio) {
            audio.pause();
            audio.removeAttribute("src");
            audio.load();
        }
        setState({
            isPlaying: false,
            currentTime: 0,
            duration: 0,
            activeId: null,
        });
    }, []);

    const seek = useCallback((ratio: number) => {
        const audio = audioRef.current;
        if (audio && audio.duration) {
            audio.currentTime = ratio * audio.duration;
        }
    }, []);

    const setVolume = useCallback(
        (volume: number) => {
            const audio = getAudio();
            audio.volume = volume;
            localStorage.setItem("uminekoVolume", volume.toString());
        },
        [getAudio],
    );

    return { state, play, stop, seek, setVolume };
}

export type AudioPlayer = ReturnType<typeof useAudioPlayer>;
