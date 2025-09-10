import Intl from '../../i18n/i18n';

interface IPlaceholder {
  message: string;
  icon: JSX.Element;
}

function Placeholder(props: IPlaceholder) {
  const { message, icon } = props;
  return (
    <div
      className="flex flex-col justify-around items-center h-96 max-w-2xl
      p-10 gap-10 relative bg-secondary rounded-lg"
    >
      {icon}
      <h1 className="mas-banner text-center cursor-default">
        {Intl.t('placeholder.teaser-banner')}
      </h1>
      <p className="mas-buttons text-center">{message}</p>
    </div>
  );
}

export default Placeholder;
