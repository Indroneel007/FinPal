import { motion } from 'framer-motion';
import Prism3D from './Prism3D';
import FeaturesSection from './FeaturesSection';

export default function LandingSection() {
  return (
    <>
    <section className="relative min-h-screen w-full overflow-x-hidden flex flex-col items-center justify-center">
      {/* Animated colorful blobs background */}
      <div className="absolute inset-0 -z-10 overflow-hidden">
        {/* Subtle animated blobs */}
        <div className="absolute top-[-12%] left-[-12%] w-[38vw] h-[38vw] bg-pink-500 opacity-40 rounded-full filter blur-3xl animate-blob1" />
        <div className="absolute bottom-[-12%] right-[-10%] w-[38vw] h-[38vw] bg-blue-500 opacity-40 rounded-full filter blur-3xl animate-blob2" />
        <div className="absolute top-[40%] left-[60%] w-[32vw] h-[32vw] bg-green-500 opacity-30 rounded-full filter blur-2xl animate-blob3" />
        <div className="absolute bottom-[20%] left-[30%] w-[35vw] h-[35vw] bg-yellow-400 opacity-30 rounded-full filter blur-2xl animate-blob4" />
        <div className="absolute top-[10%] right-[20%] w-[28vw] h-[28vw] bg-purple-600 opacity-30 rounded-full filter blur-2xl animate-blob5" />
        <div className="absolute bottom-[10%] left-[10%] w-[28vw] h-[28vw] bg-fuchsia-500 opacity-25 rounded-full filter blur-2xl animate-blob6" />
        {/* Fast animated gradient overlay, less visible */}
        <div className="absolute inset-0 bg-[linear-gradient(120deg,_#ff3e3e,_#ffb800,_#22d3ee,_#22c55e,_#6366f1,_#a21caf,_#f472b6,_#ff3e3e)] bg-[length:200%_200%] animate-gradientMoveFast opacity-40 mix-blend-lighten" />
        {/* Restore darker base overlay for chill effect */}
        <div className="absolute inset-0 bg-gradient-to-br from-[#18181b] via-[#27272a] to-[#3b0764] opacity-90" />
      </div>
      {/* Prism 3D with rainbow */}
      <div className="relative flex flex-col items-center mb-12">
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 1 }}
        >
          <Prism3D />
        </motion.div>
      </div>
      {/* One-liner */}
      <motion.h1
        initial={{ opacity: 0, y: 40 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 1.2, duration: 1 }}
        className="text-4xl md:text-6xl font-extrabold text-white text-center mb-6 drop-shadow-lg"
      >
        Take control of your finances with <span className="text-violet-400">FinPal</span>
      </motion.h1>
      <motion.p
        initial={{ opacity: 0, y: 40 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 1.6, duration: 1 }}
        className="text-lg md:text-2xl text-gray-300 text-center mb-10 max-w-xl"
      >
        Your all-in-one personal finance tracker, simplified.
      </motion.p>

      {/* Features Section */}
    </section>
    <section>
      <div className="w-full flex justify-center">
        <FeaturesSection />
      </div>
    </section>
    </>
  );
}
