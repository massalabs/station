import { Theme } from '@massalabs/react-ui-kit';
import { ReactNode, useEffect, useState } from 'react';

import { Foundation } from './Dashboard/Foundation';
import { Bridge } from './Dashboard/Bridge';
import { MassaLabs } from './Dashboard/Massalabs';
import { Explorer } from './Dashboard/Explorer';
import { Purrfect } from './Dashboard/Purrfect';
import { Dusa } from './Dashboard/Dusa';

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
      className="grid lg:grid-cols-10  grid-rows-3 gap-4 h-fit p-4"
      data-testid="dashboard-station"
    >
      <div className="col-start-1 col-span-2 row-span-3 ">
        <div className={`${sizeClass}`}>{images[0]}</div>
      </div>
      <div className="col-start-3 col-span-2 row-start-1 row-span-2">
        <Foundation />
      </div>
      <div className="col-start-5 col-span-2 row-start-1 row-span-2">
        <Bridge />
      </div>
      <div className="col-start-7 col-span-4 row-start-1 row-span-1">
        <MassaLabs />
      </div>
      <div className="col-start-3 col-span-4 row-start-3 row-span-1">
        <Explorer />
      </div>
      <div className="col-start-7 col-span-2 row-start-2 row-span-2">
        <Purrfect />
      </div>
      <div className="col-start-9 col-span-2 row-start-2 row-span-2">
        <Dusa />
      </div>
    </div>
  );
}
