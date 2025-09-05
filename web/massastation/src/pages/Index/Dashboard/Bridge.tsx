import { RedirectTile } from '@massalabs/react-ui-kit';
import bridgeOutline from '../../../assets/dashboard/BridgeOutline.svg';
import { motion } from 'framer-motion';

export function Bridge() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url="https://bridge.massa.net"
        className="bg-c-default hover:bg-c-hover rounded-lg border border-c-default h-full relative overflow-hidden"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-c-default/80 to-c-default"></div>
        <div className="relative z-10 h-full p-8">
          <div className="flex flex-col h-full justify-between gap-6">
            <motion.img
              src={bridgeOutline}
              alt="Bridge"
              className="w-14 h-14"
              whileHover={{ 
                scale: 1.1,
                rotate: 5,
                transition: { duration: 0.3 }
              }}
            />
            <div className="text-brand font-bold text-4xl leading-tight">
              <div>Massa</div>
              <div>Bridge</div>
            </div>
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
