"use client";

import { Player } from "@remotion/player";
import { ArchitectureFlow } from "../animations/ArchitectureFlow";

export default function ArchitecturePlayer() {
  return (
    <Player
      component={ArchitectureFlow}
      compositionWidth={900}
      compositionHeight={300}
      durationInFrames={180}
      fps={30}
      loop
      autoPlay
      style={{ width: "100%", height: "auto" }}
      controls={false}
    />
  );
}
