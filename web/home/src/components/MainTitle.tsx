const MainTitle = (title: { title: string }) => {
  return (
    <div className="display justify-between items-center flex-row flex text-font w-full">
      <div className="flex-row flex">
        <p className="text-brand">↳</p> {title.title}
      </div>
      <div
        onClick={() => {
          window.open('/store');
        }}
        className="justify-end center cursor-pointer"
      >
        <h2 className="button text-font underline">Manage plugins ↗</h2>
      </div>
    </div>
  );
};

export default MainTitle;
