import sys

import tensorflow as tf
from tensorflow.keras.applications.mobilenet import MobileNet
from tensorflow.keras.layers import Dense, GlobalAveragePooling2D
from tensorflow.keras.models import Model
from tensorflow.keras.optimizers import SGD
import tensorflow.keras.preprocessing as pp 

def export_access_pattern(current_epoch, epochs, paths_list):
   with open('access_pattern_'+str(current_epoch), 'w') as f:
      for paths in paths_list:
         in_epoch = paths.numpy()
         for path in in_epoch:
            f.write(path.decode('utf-8') + '\n')

if len(sys.argv) != 4:
  print('Usage: python3 image_training.py <dataset_path> <batch_size> <epochs>')
  exit(1)

dataset_path = sys.argv[1]
batch_size = int(sys.argv[2])
epochs = int(sys.argv[3])

# Load MobileNet model without top layers
base_model = MobileNet(weights='imagenet', include_top=False, input_shape=(224, 224, 3))

# Add top layers to the model
x = base_model.output
x = GlobalAveragePooling2D()(x)
x = Dense(1024, activation='relu')(x)
predictions = Dense(1000, activation='softmax')(x)
model = Model(inputs=base_model.input, outputs=predictions)

# Freeze the base model layers
for layer in base_model.layers:
    layer.trainable = False

# Compile the model
model.compile(optimizer=SGD(learning_rate=0.001), loss='categorical_crossentropy', metrics=['accuracy'])

# Load and preprocess the data
train_datagen = pp.image.ImageDataGenerator(rescale=1./255)
dataset = pp.image_dataset_from_directory(
        dataset_path,
        image_size=(224, 224),
        batch_size=batch_size,
        labels='inferred',
        label_mode='categorical',
        shuffle=False)

# dataset을 셔플하면 file_paths 속성이 사라지므로, ZipDataSet으로 변경
file_paths = tf.data.Dataset.from_tensor_slices((dataset.file_paths)).batch(batch_size)
dataset_with_path = tf.data.Dataset.zip((dataset, file_paths))

# Fine-tune the model
for i in range(epochs):
   if shuffle:
      dataset_with_path = dataset_with_path.shuffle(10000)
   print('Epoch ' + str(i+1) + '/' + str(epochs))
   image_dataset = dataset_with_path.map(lambda img, path: img).prefetch(tf.data.AUTOTUNE)
   path_dataset = tf.data.Dataset.from_tensor_slices(dataset.file_paths).batch(batch_size)
   export_access_pattern(i+1, epochs, path_dataset) # Access 순서
   model.fit(
         image_dataset,
         steps_per_epoch=10000//batch_size,
         epochs=1,
         shuffle=False)