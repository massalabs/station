import { DomainModel } from '../../../models/DomainModel';

import { Description } from '@massalabs/react-ui-kit';
import { FiGlobe } from 'react-icons/fi';

interface DomainModelItemProps {
  website: DomainModel;
}

export default function DomainModelItem(props: DomainModelItemProps) {
  const { website } = props;

  return (
    <Description
      variant="secondary"
      preIcon={<FiGlobe />}
      title={website.name}
      website={website.name + '.massa'}
      description={website.description}
    />
  );
}
