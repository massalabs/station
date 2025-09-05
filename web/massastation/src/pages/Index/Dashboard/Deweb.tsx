import { motion } from 'framer-motion';
import { useState } from 'react';
import { RedirectTile } from '@massalabs/react-ui-kit';
import { routeFor } from '@/utils/utils';
import dewebLogo from '@/assets/dashboard/deweb.svg';


export function Deweb() {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url={routeFor('deweb')}
        data-testid="deweb"
        className="w-full h-full rounded-md flex flex-col"
      >
        <div className="w-full h-full rounded-md bg-c-default relative overflow-hidden">
          <div className="absolute inset-0 bg-gradient-to-br from-c-default/80 to-c-default"></div>
          <div className="relative z-10 h-full p-6 flex items-center justify-center">
            <motion.img
              initial={false}
              animate={{
                scale: isHovered ? 1.1 : 1,
                opacity: isHovered ? 0.9 : 0.8,
                transition: { duration: 0.36 },
              }}
              src={dewebLogo}
              alt="Deweb"
              className="max-w-full max-h-full object-contain"
              onMouseEnter={() => setIsHovered(true)}
              onMouseLeave={() => setIsHovered(false)}
            />
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
