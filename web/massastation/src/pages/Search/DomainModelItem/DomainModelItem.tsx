import { DomainModel } from '@/models/DomainModel';

import { Description } from '@massalabs/react-ui-kit';
import { FiGlobe } from 'react-icons/fi';

interface DomainModelItemProps {
  website: DomainModel;
}

export default function DomainModelItem(props: DomainModelItemProps) {
  const { website } = props;

  const faviconURL = `${location.protocol + '//' + website.favicon}`;
  const url = `${location.protocol + '//' + website.name + '.massa'}`;

  return (
    <Description
      variant="secondary"
      preIcon={website.favicon ? <img src={faviconURL} /> : <FiGlobe />}
      title={website.name}
      website={website.name + '.massa'}
      description={website.description}
      onClick={() => window.open(url, '_blank')}
    />
  );
}
