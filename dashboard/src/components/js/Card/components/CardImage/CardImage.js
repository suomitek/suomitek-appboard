import React from "react";
import PropTypes from "prop-types";

const CardImage = ({ src, alt }) => (
  <div className="card-img">
    <img src={src} alt={alt} />
  </div>
);

CardImage.propTypes = {
  src: PropTypes.string.isRequired,
  // For accessibility. You can pass an empty string ("") in case you want
  // screen readers to ignore the image
  alt: PropTypes.string.isRequired,
};

export default CardImage;
