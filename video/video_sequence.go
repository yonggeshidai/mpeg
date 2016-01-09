package video

import "github.com/32bitkid/bitreader"

type sequenceHeaders struct {
	*SequenceHeader
	*SequenceExtension
}

type pictureHeaders struct {
	*GroupOfPicturesHeader
	*PictureHeader
	*PictureCodingExtension
}

type VideoSequence struct {
	bitreader.BitReader

	sequenceHeaders
	pictureHeaders

	quantisationMatricies [4]QuantisationMatrix
	frameStore
}

func NewVideoSequence(br bitreader.BitReader) VideoSequence {
	return VideoSequence{
		BitReader: br,
	}
}

func (vs *VideoSequence) sequence_extension() (err error) {
	vs.SequenceExtension, err = sequence_extension(vs)
	return
}

func (vs *VideoSequence) group_of_pictures_header() (err error) {
	vs.GroupOfPicturesHeader, err = group_of_pictures_header(vs)
	return
}

func (vs *VideoSequence) picture_header() (err error) {
	vs.PictureHeader, err = picture_header(vs)
	return
}

func (vs *VideoSequence) picture_coding_extension() (err error) {
	vs.PictureCodingExtension, err = picture_coding_extension(vs)
	return
}

type dcDctPredictors [3]int32
type dcDctPredictorResetter func()

func (pred *dcDctPredictors) createResetter(intra_dc_precision uint32) dcDctPredictorResetter {
	return func() {
		resetValue := int32(1) << (7 + intra_dc_precision)
		pred[0] = resetValue
		pred[1] = resetValue
		pred[2] = resetValue
	}
}
