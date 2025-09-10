// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { Button, Spinner } from '@massalabs/react-ui-kit';
import { useState } from 'react';
import { FiArrowUpRight, FiRefreshCw } from 'react-icons/fi';
import { PluginStates } from '../DashboardStation';
import { motion } from 'framer-motion';
import nodeManager from '../../../assets/dashboard/Node_manager.svg';
import Intl from '@/i18n/i18n';


export interface NodeManagerProps {
  state?: string;
  isLoading: boolean;
  title: string;
  status?: string;
  isUpdating: boolean;
  onClickActive: () => void;
  onClickInactive: () => void;
  onUpdateClick: () => void;
}

export interface NodeManagerPluginProps {
  title: string;
  onClickActive?: () => void;
  onClickInactive?: () => void;
  isUpdating?: boolean;
  isHovered?: boolean;
}

export function ActiveNodeManager(props: NodeManagerPluginProps) {
  const { onClickActive, isHovered } = props;

  return (
    <>
      <div className="w-full h-full rounded-t-md relative overflow-hidden bg-brand">
        <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent"></div>
        <motion.img
          initial={false}
          animate={{
            rotate: isHovered ? 180 : 0,
            transition: { duration: 0.36 },
          }}
          src={nodeManager}
          alt="Node Manager"
          className="w-full h-full p-4 brightness-0 invert"
        />
      </div>
      <div
        className="w-full h-full text-primary bg-secondary flex flex-col 
        justify-end items-start p-4 rounded-b-md"
      >
        <div className="mas-subtitle text-left text-brand font-semibold mb-4 cursor-default">
          <div>Node</div>
          <div>Manager</div>
        </div>
        <div className="w-full">
          <Button 
            onClick={onClickActive} 
            preIcon={<FiArrowUpRight />}
            customClass="bg-c-default hover:bg-c-hover text-c-primary border-c-default font-semibold w-full"
          >
            Launch
          </Button>
        </div>
      </div>
    </>
  );
}

export function UpdateNodeManager(props: NodeManagerPluginProps) {
  const { onClickActive, isUpdating, isHovered } = props;

  return (
    <>
      <div className="w-full h-full rounded-t-md relative overflow-hidden bg-brand">
        <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent"></div>
        <div className="relative z-10 p-6 h-full flex items-start">
          <motion.img
            initial={false}
            animate={{
              rotate: isHovered ? 180 : 0,
              transition: { duration: 0.36 },
            }}
            src={nodeManager}
            alt="Node Manager"
            className="w-12 h-12 brightness-0 invert"
          />
        </div>
      </div>
      <div
        className="w-full h-full text-primary bg-secondary flex flex-col 
        justify-end items-start p-4 rounded-b-md"
      >
        <div className="mas-subtitle text-left text-brand font-semibold mb-4 cursor-default">
          <div>Node</div>
          <div>Manager</div>
        </div>
        <div className="flex flex-col gap-2 w-full">
          <Button onClick={onClickActive}>
            <div className="flex items-center gap-2">
              <div className={isUpdating ? 'animate-spin' : 'none'}>
                <FiRefreshCw className="text-primary" size={18} />
              </div>
              {isUpdating ? 'Updating...' : 'Update'}
            </div>
          </Button>
          <div className="text-orange-400 px-4 mas-caption font-medium cursor-default">
            <a
              className="underline"
              href="/plugin/massa-labs/node-manager/web-app/index"
              target="_blank"
            >
              {Intl.t('modules.massa-wallet.no-update')}
            </a>
          </div>
        </div>
      </div>
    </>
  );
}

export function InactiveNodeManager(props: NodeManagerPluginProps) {
  const { onClickInactive, isHovered } = props;

  return (
    <>
      <div className="w-full h-full rounded-t-md bg-neutral border-2 border-dashed 
        border-neutral/50 relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-neutral/10 to-transparent"></div>
        <div className="relative z-10 h-full flex flex-col items-center justify-center p-4">
          <motion.img
            initial={false}
            animate={{
              rotate: isHovered ? 180 : 0,
              opacity: isHovered ? 1.0 : 0.8,
              scale: isHovered ? 1.1 : 1,
              transition: { duration: 0.36 },
            }}
            src={nodeManager}
            alt="Node Manager"
            className="w-16 h-16 mb-3 grayscale brightness-50 contrast-150"
          />
          <p className="text-center text-primary font-medium text-xs cursor-default">
            {Intl.t('modules.node-manager.baseline')}
          </p>
        </div>
      </div>
      <div className="w-full h-full text-primary bg-secondary border-2 border-neutral/20 
        flex flex-col justify-end items-start rounded-b-md relative p-4">
        <div className="absolute inset-0 bg-gradient-to-t from-neutral/5 to-transparent rounded-b-md"></div>
        <div className="mas-subtitle text-left text-brand font-semibold mb-4 relative z-10 cursor-default">
          <div>Node</div>
          <div>Manager</div>
        </div>
        <div className="w-full relative z-10">
          <Button 
            variant="primary"
            customClass="bg-c-default hover:bg-c-hover text-c-primary border-c-default font-semibold w-full"
            onClick={onClickInactive}
          >
            Install
          </Button>
        </div>
      </div>
    </>
  );
}

export function CrashedNodeManager(props: NodeManagerPluginProps) {
  const { title, isHovered } = props;

  return (
    <>
      <div className="w-full h-full rounded-t-md relative overflow-hidden bg-brand">
        <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent"></div>
        <div className="relative z-10 p-6 h-full flex items-start">
          <motion.img
            initial={false}
            animate={{
              rotate: isHovered ? 180 : 0,
              transition: { duration: 0.36 },
            }}
            src={nodeManager}
            alt="Node Manager"
            className="w-12 h-12 brightness-0 invert"
          />
        </div>
      </div>
      <div className="w-full h-full py-6 text-primary bg-secondary flex flex-col items-center rounded-b-md">
        <div className="w-4/5 px-4 py-2 mas-buttons lg:h-14 flex items-center justify-center">
          <p className="text-center text-red-400 font-medium cursor-default">
            {title} has crashed. Please restart or reinstall.
          </p>
        </div>
      </div>
    </>
  );
}

export function LoadingNodeManager(props: NodeManagerPluginProps) {
  const { title, isHovered } = props;

  return (
    <>
      <div className="w-full h-full rounded-t-md relative overflow-hidden bg-brand">
        <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent"></div>
        <div className="relative z-10 p-6 h-full flex items-start">
          <motion.img
            initial={false}
            animate={{
              rotate: isHovered ? 180 : 0,
              transition: { duration: 0.36 },
            }}
            src={nodeManager}
            alt="Node Manager"
            className="w-12 h-12 brightness-0 invert"
          />
        </div>
      </div>
      <div className="w-full text-primary bg-secondary flex flex-col items-center gap-4 py-4 rounded-b-md">
        <div className="w-4/5 mas-buttons flex items-center justify-center">
          <p className="text-center text-f-primary font-medium cursor-default">
            Installing {title}...
          </p>
        </div>
        <div className="w-4/5">
          <Button 
            disabled={true}
            customClass="w-full"
          >
            <Spinner />
          </Button>
        </div>
      </div>
    </>
  );
}

export function NodeManager(props: NodeManagerProps) {
  const {
    state,
    isLoading,
    status,
    title,
    onClickActive,
    onClickInactive,
    onUpdateClick,
    isUpdating,
  } = props;

  const [isHovered, setIsHovered] = useState(false);

  const displayNodeManager = () => {
    if (isLoading) {
      return <LoadingNodeManager title={title} isHovered={isHovered} />;
    }
    if (state === PluginStates.Updateable) {
      return (
        <UpdateNodeManager
          title={title}
          onClickActive={onUpdateClick}
          isUpdating={isUpdating}
          isHovered={isHovered}
        />
      );
    }
    if (state === PluginStates.Inactive) {
      return (
        <InactiveNodeManager
          title={title}
          onClickInactive={onClickInactive}
          isHovered={isHovered}
        />
      );
    }
    if (status && status === 'Crashed') {
      return <CrashedNodeManager title={title} isHovered={isHovered} />;
    }
    return (
      <ActiveNodeManager
        title={title}
        onClickActive={onClickActive}
        isHovered={isHovered}
      />
    );
  };

  return (
    <motion.div
      onHoverStart={() => {
        setIsHovered(true);
      }}
      onHoverEnd={() => {
        setIsHovered(false);
      }}
      whileHover={{ scale: 1.03 }}
      data-testid="node-manager"
      className="w-full h-full rounded-md flex flex-col"
    >
      {displayNodeManager()}
    </motion.div>
  );
}
