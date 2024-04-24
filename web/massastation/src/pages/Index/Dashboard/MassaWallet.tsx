// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { Button, Spinner } from '@massalabs/react-ui-kit';
import Intl from '@/i18n/i18n';
import { useEffect, useState } from 'react';
import { FiArrowUpRight, FiRefreshCw } from 'react-icons/fi';
import { WalletStates } from '../DashboardStation';
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
  const { title, onClickActive, isHovered } = props;

  useEffect(() => {
    console.log('isHovered', isHovered);
  }, [isHovered]);

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
          alt="Massa Wallet"
          className="w-full h-full p-4"
        />
      </div>
      <div className="w-full h-full text-f-primary bg-secondary flex flex-col gap-1 justify-evenly items-center pb-4">
        <div className="mas-subtitle text-center">{title}</div>
        <div className="mas-body2">{Intl.t('modules.massa-wallet.desc')}</div>
        <div className="w-3/5">
          <Button onClick={onClickActive} preIcon={<FiArrowUpRight />}>
            {Intl.t('modules.massa-wallet.launch')}
          </Button>
        </div>
      </div>
    </>
  );
}

export function Updateplugin(props: MSPluginProps) {
  const { title, onClickActive, isUpdating, isHovered } = props;

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
          alt="Massa Wallet"
          className="w-full h-full p-4"
        />
      </div>
      <div className="w-full h-full text-f-primary bg-secondary flex flex-col gap-1 pb-4 justify-evenly items-center">
        <div className="mas-subtitle text-center">{title}</div>
        <div className="flex flex-col gap-2">
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
          <div className="text-s-warning px-4 mas-caption">
            <a
              className="underline"
              href="/plugin/massa-labs/massa-wallet/web-app/index"
              target="_blank"
            >
              {Intl.t('modules.massa-wallet.click')}
            </a>{' '}
            {Intl.t('modules.massa-wallet.no-update')}
          </div>
        </div>
      </div>
    </>
  );
}

export function InactivePlugin(props: MSPluginProps) {
  const { title, onClickInactive, isHovered } = props;

  return (
    <>
      <div className={`w-full h-full rounded-t-md bg-tertiary`}>
        <motion.img
          initial={false}
          animate={{
            rotate: isHovered ? 180 : 0,
            transition: { duration: 0.36 },
          }}
          src={massaWallet}
          alt="Massa Wallet"
          className="w-full h-full p-4"
        />
      </div>
      <div className="w-full h-full text-f-primary bg-secondary flex flex-col justify-center items-center rounded-b-md">
        <div className="w-4/5 px-4 py-2 mas-buttons lg:h-14 flex items-center justify-center">
          <p className="text-center">{`${title} is not installed in your station`}</p>
        </div>
        <div className="w-4/5 px-4 py-2">
          <Button onClick={onClickInactive}>
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
          alt="Massa Wallet"
          className="w-full h-full p-4"
        />
      </div>
      <div className="w-full h-full py-6 text-f-primary bg-secondary flex flex-col items-center">
        <div className="w-4/5 px-4 py-2 mas-buttons lg:h-14 flex items-center justify-center">
          <p className="text-center">{`${title} can’t be opened. Reinstall it from the Module store.`}</p>
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
          alt="Massa Wallet"
          className="w-full h-full p-4"
        />
      </div>
      <div className="w-full text-f-primary bg-secondary flex flex-col items-center gap-4 py-4">
        <div className="w-4/5 mas-buttons flex items-center justify-center">
          <p className="text-center">{`${title} installation`}</p>
        </div>
        <div className="w-3/5">
          <Button disabled={true}>
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
    if (state === WalletStates.Updateable) {
      return (
        <Updateplugin
          title={title}
          onClickActive={onUpdateClick}
          isUpdating={isUpdating}
          isHovered={isHovered}
        />
      );
    }
    if (state === WalletStates.Inactive) {
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
