import { Theme } from '@massalabs/react-ui-kit';

import { ReactNode, useEffect, useState } from 'react';

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
      className="grid lg:grid-cols-6 grid-cols-3 grid-rows-2 gap-8"
      data-testid="dashboard-station"
    >
      <div className="col-span-2 row-span-2">
        <div className={`${sizeClass}`}>{images[0]}</div>
      </div>
      <div className="col-start-3">
        <div className={`${sizeClass}`}>{images[1]}</div>
      </div>
      <div className="col-start-3 row-start-2">
        <div className={`${sizeClass}`}>{images[2]}</div>
      </div>
      <div className="row-start-3 lg:col-span-2 lg:row-span-2 lg:col-start-4 lg:row-start-1">
        <div className={`${sizeClass}`}>{images[3]}</div>
      </div>
      <div className="col-start-1 row-start-4 lg:col-start-6 lg:row-start-1">
        <div className={`${sizeClass}`}>{images[4]}</div>
      </div>
      <div
        className="col-span-2 row-span-2 col-start-2 row-start-3
      lg:col-span-1 lg:row-span-1 lg:col-start-6 lg:row-start-2"
      >
        <div className={`${sizeClass}`}>{images[5]}</div>
      </div>
    </div>
  );
}
