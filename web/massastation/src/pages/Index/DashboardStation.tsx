import { RedirectTile, Theme } from '@massalabs/react-ui-kit';
import { ReactNode, useEffect, useState } from 'react';
import Intl from '@/i18n/i18n';
import bridge from '../../assets/dashboard/Bridge.svg';
import explorer from '../../assets/dashboard/Explorer.svg';
import foundation from '../../assets/dashboard/Foundation.svg';
import massa from '../../assets/dashboard/Massa.svg';
import purrfect from '../../assets/dashboard/Purrfect.svg';
import dusa from '../../assets/dashboard/Dusa.svg';
export interface IDashboardStationProps {
  imagesDark: ReactNode[];
  imagesLight: ReactNode[];
  components: ReactNode[];
  theme?: Theme | undefined;
}

export function DashboardStation(props: IDashboardStationProps) {
  let { imagesDark, imagesLight, components, theme } = props;

  const [images, setImages] = useState<ReactNode[]>([]);
  const sizeClass = 'h-full w-full';

  useEffect(() => {
    let imageList: ReactNode[] = [...components];

    const diff = imagesDark.length - components.length;
    for (let i = 0; i < diff; i++) {
      const imageToAdd =
        theme === 'theme-dark' ? imagesDark[i] : imagesLight[i];
      imageList.push(imageToAdd);
    }

    setImages(imageList);
  }, [theme, components, imagesDark, imagesLight]);

  return (
    <div
      className="grid lg:grid-cols-10 grid-cols-3 grid-rows-3 gap-4 h-fit p-4"
      data-testid="dashboard-station"
    >
      <div className="col-start-1 col-span-2 row-span-3 ">
        <div className={`${sizeClass}`}>{images[0]}</div>
      </div>
      <div className="col-start-3 col-span-2 row-start-1 row-span-2">
        <RedirectTile
          size="cs"
          customSize="h-full"
          url="https://bridge.massa.net"
        >
          <div className="h-full flex flex-col gap-4">
            <img width={60} height={60} src={bridge} alt="Massa Bridge" />
            <p className="mas-subtitle">Massa Bridge</p>
            <p>{Intl.t('dashboard.explorer-desc')}</p>
          </div>
        </RedirectTile>
      </div>
      {/* Massa Foundation plugin from col 5 to col 7*/}
      <div className="col-start-5 col-span-2 row-start-1 row-span-2">
        <RedirectTile
          size="cs"
          customSize="h-full"
          url="https://massa.foundation/"
        >
          <div className="h-full flex flex-col gap-4">
            <img
              width={60}
              height={60}
              src={foundation}
              alt="Massa Foundation"
            />
            <p className="mas-subtitle">Massa Foundation</p>
            <p>{Intl.t('dashboard.foundation-desc')}</p>
          </div>
        </RedirectTile>
      </div>
      {/* Massa Website plugin from col 7 to col 10*/}
      <div className="col-start-7 col-span-4 row-start-1 row-span-1">
        <RedirectTile size="cs" customSize="h-full" url="https://massa.net/">
          <div className="h-fit flex items-center gap-4">
            <img width={60} height={60} src={massa} alt="Massa Website" />
            <div className="flex flex-col gap-2">
              <p className="mas-subtitle">Massa Website</p>
              <p>{Intl.t('dashboard.massalabs-desc')}</p>
            </div>
          </div>
        </RedirectTile>
      </div>
      {/* Massa Explorer plugin from col 7 to col 10*/}
      <div className="col-start-3 col-span-4 row-start-3 row-span-1">
        <RedirectTile
          size="cs"
          customSize="h-full"
          url="https://explorer.massa.net/"
        >
          <div className="h-fit flex items-center gap-4">
            <img width={60} height={60} src={explorer} alt="Massa Explorer" />
            <div className="flex flex-col gap-2">
              <p className="mas-subtitle">Massa Explorer</p>
              <p>{Intl.t('dashboard.explorer-desc')}</p>
            </div>
          </div>
        </RedirectTile>
      </div>
      <div className="col-start-7 col-span-2 row-start-2 row-span-2">
        <RedirectTile
          size="cs"
          customSize="h-full"
          url="https://www.purrfectuniverse.com/"
        >
          <div className="h-full flex flex-col gap-4">
            <img
              width={60}
              height={60}
              src={purrfect}
              alt="Purrfect Universe"
            />
            <p className="mas-subtitle">Purrfect Universe</p>
            <p>{Intl.t('dashboard.purrfect-desc')}</p>
          </div>
        </RedirectTile>
      </div>
      <div className="col-start-9 col-span-2 row-start-2 row-span-2">
        <RedirectTile size="cs" customSize="h-full" url="https://dusa.io/">
          <div className="h-full flex flex-col gap-4">
            <img width={60} height={60} src={dusa} alt="Dusa Protocol" />
            <p className="mas-subtitle">Dusa Protocol</p>
            <p>{Intl.t('dashboard.dusa-desc')}</p>
          </div>
        </RedirectTile>
      </div>
    </div>
  );
}
