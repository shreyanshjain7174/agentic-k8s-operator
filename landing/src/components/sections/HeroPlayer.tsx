"use client";

import { Player } from "@remotion/player";
import { HeroAnimation } from "../animations/HeroAnimation";

export default function HeroPlayer() {
  return (
    <Player
      component={HeroAnimation}
      compositionWidth={800}
      compositionHeight={600}
      durationInFrames={300}
      fps={30}
      loop
      autoPlay
      style={{ width: "100%", height: "100%" }}
      controls={false}
    />
  );
}
