import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, PerspectiveCamera } from '@react-three/drei';
import { useRef, useMemo } from 'react';
import { Texture } from 'three';

function RainbowBeam({ position, rotation, texture }) {
  return (
    <mesh position={position} rotation={rotation}>
      <boxGeometry args={[2.5, 0.13, 0.13]} />
      <meshBasicMaterial color="white" transparent map={texture} />
    </mesh>
  );
}

function Prism() {
  const mesh = useRef();
  useFrame(() => {
    if (mesh.current) {
      mesh.current.rotation.y += 0.015;
      mesh.current.rotation.x = Math.sin(Date.now() * 0.0005) * 0.1;
    }
  });
  return (
    <group ref={mesh}>
      <mesh castShadow receiveShadow>
        {/* Square-based pyramid: ConeGeometry with 4 segments */}
        <coneGeometry args={[1.2, 2, 4]} />
        <meshPhysicalMaterial
          color="#f3e7fc"
          metalness={0.7}
          roughness={0.8}
          reflectivity={1}
          clearcoat={1}
          clearcoatRoughness={0.05}
          transmission={0.93}
          ior={1.48}
          thickness={1}
          specularIntensity={1}
          sheen={1}
          sheenColor="#f3e7fc"
          opacity={0.96}
          transparent
        />
      </mesh>
    </group>
  );
}

export default function Prism3D() {
  return (
    <div style={{ width: '200px', height: '400px', margin: '0 auto' }}>
      <Canvas shadows dpr={[1, 2]}>
        <ambientLight intensity={1.1} />
        <directionalLight position={[5, 5, 5]} intensity={0.8} castShadow />
        <PerspectiveCamera makeDefault position={[0, 2, 6]} fov={40} />
        <Prism />
        <OrbitControls enableZoom={false} enablePan={false} autoRotate autoRotateSpeed={-0.7} />
      </Canvas>
    </div>
  );
}
