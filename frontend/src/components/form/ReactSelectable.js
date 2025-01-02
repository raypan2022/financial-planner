import Select from 'react-select';

const ReactSelectable = ({ mandatory, name, title, value, onChange, options, errorDiv, errorMsg }) => {
  return (
    <div className="mb-3">
      <label htmlFor={name} className="form-label">
        {title}{mandatory && <span className="text-danger"> *</span>}
      </label>
      <Select
        isClearable
        classNamePrefix={"form-select"}
        value={value.value ? value : null}
        onChange={onChange}
        options={options}
        placeholder="Select an option..."
      />
      <div className={errorDiv}>{errorMsg}</div>
    </div>
  );
};

export default ReactSelectable;
