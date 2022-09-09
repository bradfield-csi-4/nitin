#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t sum1 = 0, sum2 = 0;
   long length = vec_length(u); // we can assume both vectors are same length
   long i;

   for (i = 0; i < length; i+=2) {
       sum1 += u->data[i] * v->data[i];
       sum2 += u->data[i+1] * v->data[i+1];
   }

   for (; i < length; i++) {
       sum1 += u->data[i] * v->data[i;
   }

   return sum1 + sum2;
}

data_t dotproduct_original(vec_ptr u, vec_ptr v) {
    data_t sum = 0, u_val, v_val;

    for (long i = 0; i < vec_length(u); i++) { // we can assume both vectors are same length
        get_vec_element(u, i, &u_val);
        get_vec_element(v, i, &v_val);
        sum += u_val * v_val;
    }
    return sum;
}
