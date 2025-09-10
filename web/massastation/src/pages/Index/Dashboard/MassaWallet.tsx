// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { Button, Spinner } from '@massalabs/react-ui-kit';
import Intl from '@/i18n/i18n';
import { useState } from 'react';
import { FiArrowUpRight, FiRefreshCw } from 'react-icons/fi';
import { PluginStates } from '../DashboardStation';
import { motion } from 'framer-motion';
import massaWallet from '../../../assets/dashboard/MassaWallet.svg';

export interface PluginWalletProps {
  state?: string;
  isLoading: boolean;
  title: string;
  status?: string;
  isUpdating: boolean;
  onClickActive: () => void;
  onClickInactive: () => void;
  onUpdateClick: () => void;
}

export interface MSPluginProps {
  title: string;
  onClickActive?: () => void;
  onClickInactive?: () => void;
  isUpdating?: boolean;
  isHovered?: boolean;
}

export function ActivePlugin(props: MSPluginProps) {
  const { onClickActive, isHovered } = props;

  return (
    <>
      <div className={`w-full h-full rounded-t-md bg-brand`}>
        <motion.img
          initial={false}
          animate={{
            rotate: isHovered ? 180 : 0,
            transition: { duration: 0.36 },
          }}
          src={massaWallet}
          alt={Intl.t('modules.massa-wallet.title')}
          className="w-full h-full p-4"
        />
      </div>
      <div
        className="w-full h-full text-f-primary bg-secondary flex flex-col 
        justify-end items-start p-4 rounded-b-md"
      >
        <div className="mas-subtitle text-left text-brand font-semibold mb-4 cursor-default">
          <div>Massa</div>
          <div>Wallet</div>
        </div>
        <div className="w-full">
          <Button 
            onClick={onClickActive} 
            preIcon={<FiArrowUpRight />}
            customClass="bg-c-default hover:bg-c-hover text-c-primary border-c-default font-semibold w-full"
          >
            {Intl.t('modules.massa-wallet.launch')}
          </Button>
        </div>
      </div>
    </>
  );
}

export function Updateplugin(props: MSPluginProps) {
  const { onClickActive, isUpdating, isHovered } = props;

  return (
    <>
      <div className={`w-full h-full rounded-t-md bg-brand`}>
        <motion.img
          initial={false}
          animate={{
            rotate: isHovered ? 180 : 0,
            transition: { duration: 0.36 },
          }}
          src={massaWallet}
          alt={Intl.t('modules.massa-wallet.title')}
          className="w-full h-full p-4"
        />
      </div>
      <div
        className="w-full h-full text-f-primary bg-secondary flex flex-col 
        justify-end items-start p-4 rounded-b-md"
      >
        <div className="mas-subtitle text-left text-brand font-semibold mb-4 cursor-default">
          <div>Massa</div>
          <div>Wallet</div>
        </div>
        <div className="flex flex-col gap-2 w-full">
          <Button onClick={onClickActive}>
            <div className="flex items-center gap-2">
              <div className={isUpdating ? 'animate-spin' : 'none'}>
                <FiRefreshCw className="text-primary" size={18} />
              </div>
              {isUpdating
                ? Intl.t('modules.massa-wallet.updating')
                : Intl.t('modules.massa-wallet.update')}
            </div>
          </Button>
          <div className="text-orange-400 px-4 mas-caption font-medium cursor-default">
            <a
              className="underline"
              href="/plugin/massa-labs/massa-wallet/web-app/index"
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

export function InactivePlugin(props: MSPluginProps) {
  const { onClickInactive, isHovered } = props;

  return (
    <>
      <div className={`w-full h-full rounded-t-md bg-neutral border-2 border-dashed 
        border-neutral/50 relative overflow-hidden`}>
        <div className="absolute inset-0 bg-gradient-to-br from-neutral/10 to-transparent"></div>
        <div className="relative z-10 h-full flex flex-col items-center p-4">
          <div className="flex-1 flex items-center justify-center">
            <motion.img
              initial={false}
              animate={{
                rotate: isHovered ? 180 : 0,
                opacity: isHovered ? 1.0 : 0.8,
                scale: isHovered ? 1.1 : 1,
                transition: { duration: 0.36 },
              }}
              src={massaWallet}
              alt={Intl.t('modules.massa-wallet.title')}
              className="w-16 h-16 mb-3 grayscale brightness-50 contrast-200 dark-v2:brightness-50 dark-v2:contrast-150"
            />
          </div>
          <p className="text-center text-primary font-medium text-xs cursor-default">
            {Intl.t('modules.massa-wallet.baseline')}
          </p>
        </div>
      </div>
      <div className="w-full h-full text-primary bg-secondary border-2 border-neutral/20 
        flex flex-col justify-end items-start rounded-b-md relative p-4">
        <div className="absolute inset-0 bg-gradient-to-t from-neutral/5 to-transparent rounded-b-md"></div>
        <div className="mas-subtitle text-left text-brand font-semibold mb-4 relative z-10 cursor-default">
          <div>Massa</div>
          <div>Wallet</div>
        </div>
        <div className="w-full relative z-10">
          <Button 
            variant="primary"
            customClass="bg-c-default hover:bg-c-hover text-c-primary border-c-default font-semibold w-full"
            onClick={onClickInactive}
          >
            {Intl.t('modules.massa-wallet.install')}
          </Button>
        </div>
      </div>
    </>
  );
}

export function CrashedPlugin(props: MSPluginProps) {
  const { title, isHovered } = props;

  return (
    <>
      <div className={`w-full h-full rounded-t-md bg-brand`}>
        <motion.img
          initial={false}
          animate={{
            rotate: isHovered ? 180 : 0,
            transition: { duration: 0.36 },
          }}
          src={massaWallet}
          alt={Intl.t('modules.massa-wallet.title')}
          className="w-full h-full p-4"
        />
      </div>
      <div className="w-full h-full py-6 text-primary bg-secondary flex flex-col items-center rounded-b-md">
        <div className="w-4/5 px-4 py-2 mas-buttons lg:h-14 flex items-center justify-center">
          <p className="text-center">
            {' '}
            {Intl.t('modules.massa-wallet.crash', {
              title: title,
            })}
          </p>
        </div>
      </div>
    </>
  );
}

export function LoadingPlugin(props: MSPluginProps) {
  const { title, isHovered } = props;

  return (
    <>
      <div className={`w-full h-full rounded-t-md bg-brand`}>
        <motion.img
          initial={false}
          animate={{
            rotate: isHovered ? 180 : 0,
            transition: { duration: 0.36 },
          }}
          src={massaWallet}
          alt={Intl.t('modules.massa-wallet.title')}
          className="w-full h-full p-4"
        />
      </div>
      <div className="w-full text-primary bg-secondary flex flex-col items-center gap-4 py-4 rounded-b-md">
        <div className="w-4/5 mas-buttons flex items-center justify-center">
          <p className="text-center">
            {Intl.t('modules.massa-wallet.installing', {
              title: title,
            })}
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

export function MassaWallet(props: PluginWalletProps) {
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

  const displayPlugin = () => {
    if (isLoading) {
      return <LoadingPlugin title={title} isHovered={isHovered} />;
    }
    if (state === PluginStates.Updateable) {
      return (
        <Updateplugin
          title={title}
          onClickActive={onUpdateClick}
          isUpdating={isUpdating}
          isHovered={isHovered}
        />
      );
    }
    if (state === PluginStates.Inactive) {
      return (
        <InactivePlugin
          title={title}
          onClickInactive={onClickInactive}
          isHovered={isHovered}
        />
      );
    }
    if (status && status === 'Crashed') {
      return <CrashedPlugin title={title} isHovered={isHovered} />;
    }
    return (
      <ActivePlugin
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
      data-testid="plugin-wallet"
      className="w-full h-full rounded-md flex flex-col"
    >
      {displayPlugin()}
    </motion.div>
  );
}
