import { DomainModel } from '@/models/DomainModel';

import { Description } from '@massalabs/react-ui-kit';
import { HttpStatusCode } from 'axios';
import { useEffect, useState } from 'react';
import { FiGlobe } from 'react-icons/fi';

interface DomainModelItemProps {
  website: DomainModel;
}

export default function DomainModelItem(props: DomainModelItemProps) {
  const { website } = props;

  const faviconURL = `${location.protocol + '//' + website.favicon}`;
  const url = `${location.protocol + '//' + website.name + '.massa'}`;

  const [imageDataURL, setImageDataURL] = useState('');

  useEffect(() => {
    // Try to fetch favicon if exists
    const fetchFavicon = async () => {
      try {
        const response = await fetch(faviconURL);
        if (response.status !== HttpStatusCode.Ok) {
          return;
        }
        const buffer = await response.arrayBuffer();

        if (!buffer.byteLength) {
          return;
        }
        const blob = new Blob([buffer]);
        const dataURL = URL.createObjectURL(blob);
        setImageDataURL(dataURL);
      } catch (error) {
        console.error(
          `error fetching favicon for ${website.name}.massa`,
          error,
        );
      }
    };

    fetchFavicon();
  }, [website.favicon]);

  return (
    <Description
      variant="secondary"
      preIcon={imageDataURL.length ? <img src={imageDataURL} /> : <FiGlobe />}
      title={website.name}
      website={website.name + '.massa'}
      description={website.description}
      onClick={() => window.open(url, '_blank')}
    />
  );
}
