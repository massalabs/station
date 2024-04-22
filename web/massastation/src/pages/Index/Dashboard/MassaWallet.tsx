// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { Button, Spinner } from '@massalabs/react-ui-kit';

import { ReactNode } from 'react';
import { FiArrowUpRight, FiRefreshCw } from 'react-icons/fi';
import { WalletStates } from '../DashboardStation';

export interface PluginWalletProps {
  state?: string;
  isLoading: boolean;
  title: string;
  status?: string;
  isUpdating: boolean;
  iconActive: ReactNode;
  iconInactive: ReactNode;
  onClickActive: () => void;
  onClickInactive: () => void;
  onUpdateClick: () => void;
}

export interface MSPluginProps {
  title: string;
  iconActive?: ReactNode;
  onClickActive?: () => void;
  iconInactive?: ReactNode;
  onClickInactive?: () => void;
  isUpdating?: boolean;
}

export function ActivePlugin(props: MSPluginProps) {
  const { title, iconActive, onClickActive } = props;

  return (
    <>
      <div>{iconActive}</div>
      <div className="w-full py-6 text-f-primary bg-secondary flex flex-col items-center">
        <div className="px-4 py-2 lg:h-14 mas-title text-center">
          <p className="text-xl sm:text-4xl lg:text-2xl 2xl:text-4xl">
            {title}
          </p>
        </div>
        <div className="w-4/5 px-4 py-2">
          <Button onClick={onClickActive} preIcon={<FiArrowUpRight />}>
            Launch
          </Button>
        </div>
      </div>
    </>
  );
}

export function Updateplugin(props: MSPluginProps) {
  const { title, iconActive, onClickActive, isUpdating } = props;

  return (
    <>
      <div>{iconActive}</div>
      <div className="w-full py-6 text-f-primary bg-secondary flex flex-col items-center">
        <div className="px-4 py-2 lg:h-14 mas-title text-center">
          <p className="text-xl sm:text-4xl lg:text-2xl 2xl:text-4xl">
            {title}
          </p>
        </div>
        <div className="flex flex-col gap-2 px-4 py-2">
          <Button onClick={onClickActive}>
            <div className="flex gap-2">
              {' '}
              <div className={isUpdating ? 'animate-spin' : 'none'}>
                <FiRefreshCw color={'black'} size={20} />
              </div>
              {isUpdating ? 'Updating...' : 'Update'}
            </div>
          </Button>
          <div className="text-s-warning px-4 mas-caption">
            <a
              className="underline"
              href="/plugin/massa-labs/massa-wallet/web-app/index"
              target="_blank"
            >
              Click here
            </a>{' '}
            to launch it without update
          </div>
        </div>
      </div>
    </>
  );
}

export function InactivePlugin(props: MSPluginProps) {
  const { title, iconInactive, onClickInactive } = props;

  return (
    <>
      <div>{iconInactive}</div>
      <div className="w-full h-full text-f-primary bg-secondary flex flex-col justify-center items-center rounded-b-md">
        <div className="w-4/5 px-4 py-2 mas-buttons lg:h-14 flex items-center justify-center">
          <p className="text-center">{`${title} is not installed in your station`}</p>
        </div>
        <div className="w-4/5 px-4 py-2">
          <Button onClick={onClickInactive}>Install</Button>
        </div>
      </div>
    </>
  );
}

export function CrashedPlugin(props: MSPluginProps) {
  const { title, iconActive } = props;

  return (
    <>
      {iconActive}
      <div className="w-full py-6 text-f-primary bg-secondary flex flex-col items-center">
        <div className="w-4/5 px-4 py-2 mas-buttons lg:h-14 flex items-center justify-center">
          <p className="text-center">{`${title} can’t be opened. Reinstall it from the Module store.`}</p>
        </div>
      </div>
    </>
  );
}

export function LoadingPlugin(props: MSPluginProps) {
  const { title, iconInactive } = props;

  return (
    <>
      {iconInactive}
      <div className="w-full py-6 text-f-primary bg-secondary flex flex-col items-center">
        <div className="w-4/5 px-4 py-2 mas-buttons lg:h-14 flex items-center justify-center">
          <p className="text-center">{`${title} installation`}</p>
        </div>
        <div className="w-4/5 px-4 py-2">
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
    iconActive,
    iconInactive,
    onClickActive,
    onClickInactive,
    onUpdateClick,
    isUpdating,
  } = props;

  const displayPlugin = () => {
    if (isLoading) {
      return <LoadingPlugin title={title} iconInactive={iconInactive} />;
    }
    if (state === WalletStates.Updateable) {
      return (
        <Updateplugin
          title={title}
          iconActive={iconActive}
          onClickActive={onUpdateClick}
          isUpdating={isUpdating}
        />
      );
    }
    if (state === WalletStates.Inactive) {
      return (
        <InactivePlugin
          title={title}
          iconInactive={iconInactive}
          onClickInactive={onClickInactive}
        />
      );
    }
    if (status && status === 'Crashed') {
      return <CrashedPlugin title={title} iconActive={iconActive} />;
    }
    return (
      <ActivePlugin
        title={title}
        iconActive={iconActive}
        onClickActive={onClickActive}
      />
    );
  };

  return (
    <div
      data-testid="plugin-wallet"
      className="w-full h-full rounded-md flex flex-col"
    >
      {displayPlugin()}
    </div>
  );
}
